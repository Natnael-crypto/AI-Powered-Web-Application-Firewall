// cypress/e2e/topThreatsChart.cy.js

describe('Top Threats Chart - Dashboard', () => {
  beforeEach(() => {
    localStorage.setItem('token', 'fake-valid-token')

    cy.intercept('GET', '/api/auth/me', {fixture: 'user.json'})
    cy.intercept('GET', '/api/statistics/top-threats*', {
      fixture: 'topThreats.json',
    }).as('getTopThreats')

    cy.visit('/dashboard')
    cy.wait('@getTopThreats')
  })

  it('renders pie chart for top threats', () => {
    cy.contains('Top 5 Threat Types').should('exist')
    cy.get('.recharts-pie-sector').should('have.length.at.least', 1)
  })

  it('shows tooltip on segment hover', () => {
    cy.get('.recharts-pie-sector').first().trigger('mouseover')
    cy.get('.recharts-tooltip-wrapper').should('be.visible')
  })
})
