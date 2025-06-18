describe('Attack Log Table', () => {
  beforeEach(() => {
    localStorage.setItem('token', 'fake-valid-token')
    cy.intercept('GET', '/api/requests*', {fixture: 'requests.json'}).as('getRequests')
    cy.intercept('GET', '/api/is-logged-in', {fixture: 'user.json'}).as('getUser')
    cy.visit('/attacks-logs')
    cy.wait(['@getUser', '@getRequests'])
  })

  it('should display the table with request logs', () => {
    // Verify table headers
    cy.contains('Status').should('be.visible')
    cy.contains('Application').should('be.visible')
    cy.contains('Method').should('be.visible')
    cy.contains('URL').should('be.visible')
    cy.contains('Threat Type').should('be.visible')
    cy.contains('IP').should('be.visible')
    cy.contains('Location').should('be.visible')
    cy.contains('Code').should('be.visible')
    cy.contains('Time').should('be.visible')

    // Verify at least one row is displayed
    cy.get('tbody tr').should('have.length.at.least', 1)
  })

  it('should show loading spinner while fetching data', () => {
    // Override the request with a delayed version to simulate loading
    cy.intercept('GET', '/api/requests*', {
      delay: 1000,
      fixture: 'requests.json',
    }).as('slowRequest')

    // Reload page to trigger the new delayed request
    cy.visit('/attacks-logs')
    cy.get('[data-testid="loading-spinner"]').should('be.visible')
    cy.wait('@slowRequest')
    cy.get('[data-testid="loading-spinner"]').should('not.exist')
  })

  it('should open request details modal when a row is clicked', () => {
    // Wait for rows to be present before clicking
    cy.get('tbody tr').should('have.length.at.least', 1)
    cy.get('tbody tr').first().click({force: true})

    // Verify modal opens with request details
    cy.get('[role="dialog"]').should('be.visible')
    cy.contains('Request Details').should('be.visible')

    // Close the modal
    cy.get('button').contains('Close').click({force: true})
    cy.get('[role="dialog"]').should('not.exist')
  })

  it('should paginate through the table', () => {
    // Mock second page response
    cy.intercept('GET', '/api/requests?page=2', {
      fixture: 'requests-page2.json',
    }).as('getPage2')

    // Click next button
    cy.get('button').contains('Next').click({force: true})
    cy.wait('@getPage2')
    cy.contains('Page 2 of').should('be.visible')

    // Click previous button
    cy.get('button').contains('Prev').click({force: true})
    cy.contains('Page 1 of').should('be.visible')
  })

  it('should download CSV when Generate Request button is clicked', () => {
    // Mock CSV download response
    cy.intercept('GET', '/api/generate-csv*', {
      fixture: 'requests.csv',
      headers: {
        'Content-Type': 'text/csv',
        'Content-Disposition': 'attachment; filename=logs.csv',
      },
    }).as('downloadCSV')

    // Stub the download process
    cy.window().then(win => {
      cy.stub(win.URL, 'createObjectURL').returns('blob:mock-url')
      cy.stub(win.URL, 'revokeObjectURL')
    })

    // Click download button
    cy.get('#generateRequest').click({force: true})
    cy.wait('@downloadCSV')

    // Verify download was triggered
    cy.window().its('URL.createObjectURL').should('be.called')
  })
})
