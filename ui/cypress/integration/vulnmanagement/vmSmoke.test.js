import { url, selectors } from '../../constants/VulnManagementPage';
import withAuth from '../../helpers/basicAuth';

describe('Smoke test for vulnmanagement', () => {
    withAuth();
    it('VulnManagement tile link is present and lands  on dashboard page', () => {
        cy.visit('/main/dashboard');
        cy.get(selectors.vmDBPageTileLink).click({ force: true });
        cy.url().should('contain', url.dashboard);
    });
});
