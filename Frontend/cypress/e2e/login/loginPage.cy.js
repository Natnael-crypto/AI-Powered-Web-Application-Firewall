describe('Login Page', () => {
  beforeEach(() => {
    cy.visit('/login')
  })

  it('renders the login form correctly', () => {
    cy.contains('Welcome Back').should('be.visible')
    cy.get('input[placeholder="Username"]').should('exist')
    cy.get('input[placeholder="Password"]').should('exist')
    cy.get('button').contains('Sign In').should('be.visible')
  })

  it('shows validation errors when fields are empty', () => {
    cy.get('button').contains('Sign In').click()
    cy.contains('Username is required').should('be.visible')
    cy.contains('Password is required').should('be.visible')
  })

  it('shows error on invalid credentials', () => {
    cy.intercept('POST', '/api/login', {
      statusCode: 401,
      body: {message: 'Invalid username or password'},
    }).as('loginFail')

    cy.get('input[placeholder="Username"]').type('invalidUser')
    cy.get('input[placeholder="Password"]').type('wrongPassword')
    cy.get('button').contains('Sign In').click()

    cy.wait('@loginFail')
    cy.contains('Invalid username or password').should('be.visible')
  })

  it('logs in successfully and redirects to dashboard', () => {
    cy.intercept('POST', '/api/login', {
      statusCode: 200,
      body: {token: 'fake-token'},
    }).as('loginSuccess')

    cy.get('input[placeholder="Username"]').type('validUser')
    cy.get('input[placeholder="Password"]').type('correctPassword')
    cy.get('button').contains('Sign In').click()

    cy.wait('@loginSuccess')
    cy.url().should('include', '/dashboard')
    cy.then(() => {
      expect(localStorage.getItem('token')).to.eq('fake-token')
    })
  })

  it('disables the login button while loading', () => {
    cy.intercept('POST', '/api/login', req => {
      return new Promise(resolve =>
        setTimeout(() => resolve({statusCode: 200, body: {token: 'fake-token'}}), 2000),
      )
    }).as('loginDelay')

    cy.get('input[placeholder="Username"]').type('validUser')
    cy.get('input[placeholder="Password"]').type('correctPassword')
    cy.get('button').contains('Sign In').click().should('be.disabled')
    cy.wait('@loginDelay')
  })
})
