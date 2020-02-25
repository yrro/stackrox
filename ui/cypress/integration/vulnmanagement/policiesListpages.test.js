import withAuth from '../../helpers/basicAuth';
import checkFeatureFlag from '../../helpers/features';
import { url, selectors } from '../../constants/VulnManagementPage';
import { hasExpectedHeaderColumns, allChecksForEntities } from '../../helpers/vmWorkflowUtils';

describe.skip('Policies list Page and its entity detail page , related entities sub list  validations ', () => {
    before(function beforeHook() {
        // skip the whole suite if vuln mgmt isn't enabled
        if (checkFeatureFlag('ROX_VULN_MGMT_UI', false)) {
            this.skip();
        }
    });

    withAuth();

    it('should display all the columns and links expected in clusters list page', () => {
        cy.visit(url.list.policies);
        hasExpectedHeaderColumns([
            'Policy',
            'Description',
            'Policy Status',
            'Last Updated',
            'Latest Violation',
            'Severity',
            'Deployments',
            // 'Lifecycle',
            'Enforcement'
        ]);
        cy.get(selectors.tableBodyColumn).each($el => {
            const columnValue = $el.text().toLowerCase();
            if (columnValue !== 'no deployments' && columnValue.includes('deployment'))
                allChecksForEntities(url.list.policies, 'deployment');
        });
    });
});
