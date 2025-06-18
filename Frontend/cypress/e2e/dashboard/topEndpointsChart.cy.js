// cypress/e2e/topEndpointsChart.cy.js

describe('Top Endpoints Chart - Dashboard', () => {
  beforeEach(() => {
    localStorage.setItem('token', 'fake-valid-token')

    cy.intercept('GET', '/api/auth/me', {fixture: 'user.json'})
    cy.intercept('GET', '/api/statistics/top-endpoints*', {
      fixture: 'topEndpoints.json',
    }).as('getTopEndpoints')

    cy.visit('/dashboard')
    cy.wait('@getTopEndpoints')
  })

  it('renders bar chart for endpoints', () => {
    cy.contains('Top 5 Targeted Endpoints').should('exist')
    cy.get('.recharts-bar').should('have.length.at.least', 1)
  })

  it('shows tooltip on hover', () => {
    cy.get('.recharts-bar-rectangle').first().trigger('mouseover')
    cy.get('.recharts-tooltip-wrapper').should('be.visible')
  })
})
