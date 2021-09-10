import React from 'react';
import { Flex, FlexItem, Button, Divider, PageSection, Title, Badge } from '@patternfly/react-core';
import { TableComposable, Thead, Tbody, Tr, Th, Td } from '@patternfly/react-table';
import { useParams, Link } from 'react-router-dom';
import resolvePath from 'object-resolve-path';

import useTableSelection from 'hooks/useTableSelection';
import usePageState from '../hooks/usePageState';
import { Integration, getIsAPIToken, getIsClusterInitBundle } from '../utils/integrationUtils';
import tableColumnDescriptor from '../utils/tableColumnDescriptor';
import DownloadCAConfigBundle from '../ClusterInitBundles/DownloadCAConfigBundle';

type TableCellProps = {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    row: Integration;
    column: {
        Header: string;
        accessor: (data) => string | string;
    };
};

function TableCellValue({ row, column }: TableCellProps): React.ReactElement {
    let value: string;
    if (typeof column.accessor === 'function') {
        value = column.accessor(row).toString();
    } else {
        value = resolvePath(row, column.accessor).toString();
    }
    return <div>{value || '-'}</div>;
}

type IntegrationsTableProps = {
    title: string;
    integrations: Integration[];
    hasMultipleDelete: boolean;
    onDeleteIntegrations: (integration) => void;
};

function IntegrationsTable({
    title,
    integrations,
    hasMultipleDelete,
    onDeleteIntegrations,
}: IntegrationsTableProps): React.ReactElement {
    const { source, type } = useParams();
    const { getPathToCreate, getPathToEdit, getPathToViewDetails } = usePageState();
    const columns = [...tableColumnDescriptor[source][type]];
    const {
        selected,
        allRowsSelected,
        hasSelections,
        onSelect,
        onSelectAll,
        getSelectedIds,
    } = useTableSelection<Integration>(integrations);

    const isAPIToken = getIsAPIToken(source, type);
    const isClusterInitBundle = getIsClusterInitBundle(source, type);

    function onDeleteIntegrationHandler() {
        const ids = getSelectedIds();
        onDeleteIntegrations(ids);
    }

    return (
        <>
            <Flex className="pf-u-p-md">
                <FlexItem
                    spacer={{ default: 'spacerMd' }}
                    alignSelf={{ default: 'alignSelfCenter' }}
                >
                    <Flex alignSelf={{ default: 'alignSelfCenter' }}>
                        <FlexItem
                            spacer={{ default: 'spacerMd' }}
                            alignSelf={{ default: 'alignSelfCenter' }}
                        >
                            <Title headingLevel="h2" className="pf-u-color-100 pf-u-ml-sm">
                                {title}
                            </Title>
                        </FlexItem>
                        <FlexItem
                            spacer={{ default: 'spacerMd' }}
                            alignSelf={{ default: 'alignSelfCenter' }}
                        >
                            <Badge isRead>{integrations.length}</Badge>
                        </FlexItem>
                    </Flex>
                </FlexItem>
                <FlexItem align={{ default: 'alignRight' }}>
                    <Flex>
                        {hasSelections && hasMultipleDelete && (
                            <FlexItem spacer={{ default: 'spacerMd' }}>
                                <Button variant="danger" onClick={onDeleteIntegrationHandler}>
                                    Delete integrations
                                </Button>
                            </FlexItem>
                        )}
                        {isClusterInitBundle && (
                            <FlexItem spacer={{ default: 'spacerMd' }}>
                                <DownloadCAConfigBundle />
                            </FlexItem>
                        )}
                        <FlexItem spacer={{ default: 'spacerMd' }}>
                            <Link to={getPathToCreate(source, type)}>
                                <Button data-testid="add-integration">New integration</Button>
                            </Link>
                        </FlexItem>
                    </Flex>
                </FlexItem>
            </Flex>
            <Divider component="div" />
            <PageSection isFilled padding={{ default: 'noPadding' }} hasOverflowScroll>
                <TableComposable variant="compact" isStickyHeader>
                    <Thead>
                        <Tr>
                            {hasMultipleDelete && (
                                <Th
                                    select={{
                                        onSelect: onSelectAll,
                                        isSelected: allRowsSelected,
                                    }}
                                />
                            )}
                            {columns.map((column) => {
                                return (
                                    <Th key={column.Header} modifier="wrap">
                                        {column.Header}
                                    </Th>
                                );
                            })}
                            <Th aria-label="Row actions" />
                        </Tr>
                    </Thead>
                    <Tbody>
                        {integrations.map((integration, rowIndex) => {
                            const { id } = integration;
                            const actionItems = [
                                {
                                    title: (
                                        <Link to={getPathToEdit(source, type, id)}>
                                            Edit integration
                                        </Link>
                                    ),
                                    isHidden: isAPIToken || isClusterInitBundle,
                                },
                                {
                                    title: (
                                        <div className="pf-u-danger-color-100">
                                            Delete Integration
                                        </div>
                                    ),
                                    onClick: () => onDeleteIntegrations([integration.id]),
                                },
                            ].filter((actionItem) => {
                                return !actionItem?.isHidden;
                            });
                            return (
                                <Tr key={integration.id}>
                                    {hasMultipleDelete && (
                                        <Td
                                            key={integration.id}
                                            select={{
                                                rowIndex,
                                                onSelect,
                                                isSelected: selected[rowIndex],
                                            }}
                                        />
                                    )}
                                    {columns.map((column) => {
                                        if (column.Header === 'Name') {
                                            return (
                                                <Td>
                                                    <Link
                                                        to={getPathToViewDetails(source, type, id)}
                                                    >
                                                        <TableCellValue
                                                            row={integration}
                                                            column={column}
                                                        />
                                                    </Link>
                                                </Td>
                                            );
                                        }
                                        return (
                                            <Td>
                                                <TableCellValue row={integration} column={column} />
                                            </Td>
                                        );
                                    })}
                                    <Td
                                        actions={{
                                            items: actionItems,
                                        }}
                                        className="pf-u-text-align-right"
                                    />
                                </Tr>
                            );
                        })}
                    </Tbody>
                </TableComposable>
            </PageSection>
        </>
    );
}

export default IntegrationsTable;
