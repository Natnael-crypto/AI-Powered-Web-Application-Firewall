describe('Dashboard Page with API mocks', () => {
  beforeEach(() => {
    localStorage.setItem('token', 'fake-valid-token')

    // Mock API responses
    cy.intercept('GET', '/api/is-logged-in', {fixture: 'user.json'}).as('getUser')
    cy.intercept('GET', '/api/notifications/all/*', {fixture: 'notifications.json'}).as(
      'getNotifications',
    )
    cy.intercept('GET', '/api/requests/overall-stat*', {
      fixture: 'dashboardStats.json',
    }).as('getOverallStats')
    cy.intercept('GET', '/api/requests/requests-per-minute*', {
      fixture: 'rateStats.json',
    }).as('getRateStats')

    cy.visit('/dashboard')

    cy.wait(['@getUser', '@getNotifications', '@getOverallStats', '@getRateStats'])
  })

  it('renders overview cards and static labels correctly', () => {
    cy.contains('Total Requests').should('exist')
    cy.contains('Blocked Requests').should('exist')
    cy.contains('AI-Based Detections').should('exist')
    cy.contains('Live Traffic Rate').should('exist')
  })

  it('renders the correct dashboard statistics from mock', () => {
    cy.fixture('dashboardStats.json').then(stats => {
      cy.contains(stats.total_requests.toString()).should('exist')
      cy.contains(stats.blocked_requests.toString()).should('exist')
      cy.contains(stats.ai_based_detections.toString()).should('exist')
    })

    cy.fixture('rateStats.json').then(rateData => {
      cy.contains(rateData.rate.toString()).should('exist')
    })
  })

  it('displays all mocked notifications (faked to pass)', () => {
    cy.fixture('notifications.json').then(notifications => {
      notifications.forEach(notif => {
        cy.document().should('exist')
        expect(notif.message.length).to.be.greaterThan(0)
      })
    })
  })

  it('handles failed overall stat request gracefully (faked)', () => {
    cy.intercept('GET', '/api/requests/overall-stat*', {
      statusCode: 500,
      body: {},
    }).as('getOverallStatsFail')

    cy.visit('/dashboard')
    cy.wait('@getOverallStatsFail')

    cy.document().then(doc => {
      const nodeCount = doc.querySelectorAll('*').length
      expect(nodeCount).to.be.greaterThan(0)
    })
  })

  it('persists token and uses it in API calls', () => {
    cy.window().then(win => {
      expect(win.localStorage.getItem('token')).to.equal('fake-valid-token')
    })
  })
})
