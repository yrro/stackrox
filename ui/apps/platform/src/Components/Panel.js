import React from 'react';
import PropTypes from 'prop-types';

import { Tooltip, TooltipOverlay } from '@stackrox/ui-components';
import CloseButton from './CloseButton';

/*
 * PageBody is sibling of PageHeader and parent of main and side panels.
 */
export function PageBody({ children }) {
    return <div className="flex flex-1 h-full relative z-0">{children}</div>;
}

/*
 * PanelNew is parent of PanelHead and PanelBody.
 */
export function PanelNew({ children, testid }) {
    return (
        <div className="flex flex-col h-full w-full" data-testid={testid}>
            {children}
        </div>
    );
}

const defaultBorderColor = 'border-base-400';

export function getSidePanelHeadBorderColor(isDarkMode) {
    return isDarkMode ? defaultBorderColor : 'border-base-100';
}

/*
 * PanelHead is parent of the following:
 * PanelTitle or entity-specific component like EntityBreadCrumbs
 * PanelHeadEnd, which has flex end alignment
 */
export function PanelHead({ children, isDarkMode, isSidePanel = false }) {
    const borderColor = isSidePanel ? getSidePanelHeadBorderColor(isDarkMode) : defaultBorderColor;
    return <div className={`${borderColor} border-b flex h-14 w-full`}>{children}</div>;
}

export function PanelTitle({ isUpperCase, testid, text }) {
    return (
        <div
            className={`flex font-700 items-center leading-normal overflow-hidden px-4 text-base-600 tracking-wide ${
                isUpperCase ? 'uppercase' : 'capitalize'
            }`}
            data-testid={testid}
        >
            <Tooltip content={<TooltipOverlay>{text}</TooltipOverlay>}>
                <div className="line-clamp break-all">{text}</div>
            </Tooltip>
        </div>
    );
}

/*
 * PanelHeadEnd, which has flex end alignment, is parent of reusable components.
 * main panel: search filter, table pagination
 * side panel: external link button, close button
 */
export function PanelHeadEnd({ children }) {
    return <div className="flex flex-1 items-center justify-end pl-3 relative">{children}</div>;
}

/*
 * PanelBody is parent of one or more entity-specific components.
 */
export function PanelBody({ children }) {
    return <div className="h-full overflow-y-auto">{children}</div>;
}

export const headerClassName = 'flex w-full min-h-14 border-b border-base-400';

const Panel = (props) => (
    <div className={`flex flex-col w-full ${props.className} h-full`} data-testid={props.id}>
        <div className="flex-nowrap">
            <div className={props.headerClassName}>
                {props.leftButtons && (
                    <div className="flex items-center pr-3 relative border-base-400 border-r hover:bg-primary-300 hover:border-primary-300">
                        {props.leftButtons}
                    </div>
                )}
                {props.headerTextComponent ? (
                    <div className="flex" data-testid={`${props.id}-header`}>
                        {props.headerTextComponent}
                    </div>
                ) : (
                    <div
                        className={`overflow-hidden mx-4 flex text-base-600 items-center tracking-wide leading-normal font-700 ${
                            props.isUpperCase ? 'uppercase' : 'capitalize'
                        }`}
                        data-testid={`${props.id}-header`}
                    >
                        <Tooltip content={<TooltipOverlay>{props.header}</TooltipOverlay>}>
                            <div className="line-clamp break-all">{props.header}</div>
                        </Tooltip>
                    </div>
                )}

                <div
                    className={`flex items-center justify-end relative flex-1 ${
                        props.onClose ? 'pl-3' : 'px-3'
                    }`}
                >
                    {props.headerComponents && props.headerComponents}
                    {props.onClose && (
                        <CloseButton
                            onClose={props.onClose}
                            className={props.closeButtonClassName}
                            iconColor={props.closeButtonIconColor}
                        />
                    )}
                </div>
            </div>
        </div>
        <div className={`h-full overflow-y-auto ${props.bodyClassName}`}>{props.children}</div>
    </div>
);

Panel.propTypes = {
    id: PropTypes.string,
    header: PropTypes.string,
    headerTextComponent: PropTypes.element,
    headerClassName: PropTypes.string,
    bodyClassName: PropTypes.string,
    className: PropTypes.string,
    children: PropTypes.node.isRequired,
    onClose: PropTypes.func,
    closeButtonClassName: PropTypes.string,
    closeButtonIconColor: PropTypes.string,
    headerComponents: PropTypes.element,
    leftButtons: PropTypes.node,
    isUpperCase: PropTypes.bool,
};

Panel.defaultProps = {
    id: 'panel',
    header: ' ',
    headerTextComponent: null,
    headerClassName,
    bodyClassName: '',
    className: '',
    onClose: null,
    closeButtonClassName: 'border-base-400 border-l',
    closeButtonIconColor: '',
    headerComponents: null,
    leftButtons: null,
    isUpperCase: true,
};

export default Panel;
