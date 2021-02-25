import React from 'react';
import { withRouter, Link } from 'react-router-dom';
import { ExternalLink } from 'react-feather';
import onClickOutside from 'react-onclickoutside';

import CloseButton from 'Components/CloseButton';
import {
    getSidePanelHeadBorderColor,
    PanelNew,
    PanelBody,
    PanelHead,
    PanelHeadEnd,
} from 'Components/Panel';
import EntityBreadCrumbs from 'Containers/BreadCrumbs/EntityBreadCrumbs';
import { useTheme } from 'Containers/ThemeProvider';
import workflowStateContext from 'Containers/workflowStateContext';
import parseURL from 'utils/URLParser';

const WorkflowSidePanel = ({ history, location, children }) => {
    const { isDarkMode } = useTheme();
    const workflowState = parseURL(location);
    const pageStack = workflowState.getPageStack();
    const breadCrumbEntities = workflowState.stateStack.slice(pageStack.length);

    function onClose() {
        const url = workflowState.removeSidePanelParams().toUrl();
        history.push(url);
    }

    WorkflowSidePanel.handleClickOutside = () => {
        const btn = document.getElementById('panel-close-button');
        if (btn) {
            btn.click();
        }
    };

    const url = workflowState.getSkimmedStack().toUrl();
    const borderColor = getSidePanelHeadBorderColor(isDarkMode);
    const externalLink = (
        <div className="flex items-center h-full hover:bg-base-300">
            <Link
                to={url}
                data-testid="external-link"
                className={`${borderColor} border-l h-full p-4`}
            >
                <ExternalLink className="h-6 w-6 text-base-600" />
            </Link>
        </div>
    );

    return (
        <workflowStateContext.Provider value={workflowState}>
            <PanelNew testid="side-panel">
                <PanelHead isDarkMode={isDarkMode} isSidePanel>
                    <EntityBreadCrumbs workflowEntities={breadCrumbEntities} />
                    <PanelHeadEnd>
                        {externalLink}
                        <CloseButton onClose={onClose} className={`${borderColor} border-l`} />
                    </PanelHeadEnd>
                </PanelHead>
                <PanelBody>{children}</PanelBody>
            </PanelNew>
        </workflowStateContext.Provider>
    );
};

const clickOutsideConfig = {
    handleClickOutside: () => WorkflowSidePanel.handleClickOutside,
};

/*
 * If more than one SidePanel is rendered, this Pure Functional Component will need to be converted to
 * a Class Component in order to work correctly. See https://github.com/stackrox/rox/pull/3090#pullrequestreview-274948849
 */
export default onClickOutside(withRouter(WorkflowSidePanel), clickOutsideConfig);
