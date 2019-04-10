import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { types as deploymentTypes } from 'reducers/deployments';
import { actions as pageActions } from 'reducers/network/page';
import { selectors } from 'reducers';
import { sortValue } from 'sorters/sorters';

import Panel from 'Components/Panel';
import Loader from 'Components/Loader';
import TablePagination from 'Components/TablePagination';
import NoResultsMessage from 'Components/NoResultsMessage';
import Table, { rtTrActionsClassName } from 'Components/Table';
import Tooltip from 'rc-tooltip';
import * as Icon from 'react-feather';

import wizardStages from '../wizardStages';

class NamespaceDetails extends Component {
    static propTypes = {
        wizardOpen: PropTypes.bool.isRequired,
        wizardStage: PropTypes.string.isRequired,
        isFetchingNamespace: PropTypes.bool,
        onClose: PropTypes.func.isRequired,
        namespace: PropTypes.shape({}),
        networkGraphRef: PropTypes.shape({
            setSelectedNode: PropTypes.func,
            selectedNode: PropTypes.shape({}),
            onNodeClick: PropTypes.func,
            getNodeData: PropTypes.func
        })
    };

    static defaultProps = {
        namespace: {},
        isFetchingNamespace: false,
        networkGraphRef: null
    };

    constructor(props) {
        super(props);
        this.state = {
            page: 0,
            selectedNode: null
        };
    }

    componentWillReceiveProps = () => {
        this.setState({ selectedNode: null });
    };

    highlightNode = ({ data }) => {
        const { networkGraphRef } = this.props;
        if (data) {
            networkGraphRef.setSelectedNode(data);
            this.setState({ selectedNode: data });
        }
    };

    navigate = ({ data }) => () => {
        const { onNodeClick } = this.props.networkGraphRef;
        if (data) {
            onNodeClick(data);
        }
    };

    setTablePage = newPage => {
        this.setState({ page: newPage });
    };

    renderRowActionButtons = node => {
        const enableIconColor = 'text-primary-600';
        const enableIconHoverColor = 'text-primary-700';
        return (
            <div className="border-2 border-r-2 border-base-400 bg-base-100 flex">
                <Tooltip
                    placement="left"
                    mouseLeaveDelay={0}
                    overlay={<div>Navigate to Deployment</div>}
                    overlayClassName="pointer-events-none"
                >
                    <button
                        type="button"
                        className={`p-1 px-4 hover:bg-primary-200 ${enableIconColor} hover:${enableIconHoverColor}`}
                        onClick={this.navigate(node)}
                    >
                        <Icon.ArrowUpRight className="mt-1 h-4 w-4" />
                    </button>
                </Tooltip>
            </div>
        );
    };

    renderTable() {
        const columns = [
            {
                Header: 'Deployment',
                accessor: 'data.name',
                Cell: ({ value }) => <span>{value}</span>
            },
            {
                Header: 'Network Flows',
                accessor: 'data.edges',
                Cell: ({ value }) => <span>{value.length}</span>,
                sortMethod: sortValue
            },
            {
                accessor: '',
                headerClassName: 'hidden',
                className: rtTrActionsClassName,
                Cell: ({ original }) => this.renderRowActionButtons(original)
            }
        ];

        const { namespace } = this.props;
        const rows = namespace.deployments;
        if (!rows.length) return <NoResultsMessage message="No namespace deployments" />;
        return (
            <Table
                rows={rows}
                columns={columns}
                onRowClick={this.highlightNode}
                noDataText="No namespace deployments"
                page={this.state.page}
                idAttribute="data.id"
                selectedRowId={this.state.selectedNode && this.state.selectedNode.id}
            />
        );
    }

    render() {
        const { namespace, wizardOpen, wizardStage, isFetchingNamespace, onClose } = this.props;
        if (!wizardOpen || wizardStage !== wizardStages.namespaceDetails) {
            return null;
        }
        const paginationComponent = (
            <TablePagination
                page={this.state.page}
                dataLength={namespace && namespace.deployments && namespace.deployments.length}
                setPage={this.setTablePage}
            />
        );
        const subHeaderText = `${namespace.deployments.length} Deployment${
            namespace.deployments.length === 1 ? '' : 's'
        }`;
        const content = isFetchingNamespace ? <Loader /> : <div>{this.renderTable()}</div>;

        return (
            <Panel header={namespace.id} onClose={onClose}>
                <Panel
                    header={subHeaderText}
                    headerComponents={paginationComponent}
                    isUpperCase={false}
                    className="w-full h-full bg-base-100"
                >
                    <div className="w-full h-full">{content}</div>
                </Panel>
            </Panel>
        );
    }
}

const mapStateToProps = createStructuredSelector({
    wizardOpen: selectors.getNetworkWizardOpen,
    wizardStage: selectors.getNetworkWizardStage,
    namespace: selectors.getSelectedNamespace,
    isFetchingNamespace: state =>
        selectors.getLoadingStatus(state, deploymentTypes.FETCH_DEPLOYMENTS),
    networkGraphRef: selectors.getNetworkGraphRef
});

const mapDispatchToProps = {
    onClose: pageActions.closeNetworkWizard
};

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NamespaceDetails);
