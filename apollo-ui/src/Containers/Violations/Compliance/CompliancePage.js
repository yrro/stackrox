import React, { Component } from 'react';
import Table from 'Components/Table';
import Tabs from 'Components/Tabs';
import TabContent from 'Components/TabContent';
import Select from 'Components/Select';

import ComplianceBenchmarksSidePanel from 'Containers/Violations/Compliance/ComplianceBenchmarksSidePanel';

import axios from 'axios';
import emitter from 'emitter';
import queryString from 'query-string';
import dateFns from 'date-fns';
import { setTimeout } from 'timers';
import { ClipLoader } from 'react-spinners';

class CompliancePage extends Component {
    constructor(props) {
        super(props);

        this.params = {};

        this.timeout = null;

        this.state = {
            tab: {
                headers: [{ text: 'CIS Docker Benchmark', disabled: false }, { text: 'Swarm Benchmark', disabled: true }, { text: 'Kubernetes Benchmark', disabled: true }]
            },
            category: {
                options: ['No scheduled scanning', 'Scan every 24 hours', 'Scan every 2 days', 'Scan every week']
            },
            table: {
                columns: [
                    { key: 'definition.name', label: 'Name' },
                    { key: 'definition.description', label: 'Description' },
                    {
                        key: 'aggregatedResults.PASS', label: 'Pass', default: 0, align: 'right'
                    },
                    {
                        key: 'aggregatedResults.INFO', label: 'Info', default: 0, align: 'right'
                    },
                    {
                        key: 'aggregatedResults.WARN', label: 'Warn', default: 0, align: 'right'
                    },
                    {
                        key: 'aggregatedResults.NOTE', label: 'Note', default: 0, align: 'right'
                    }
                ],
                rows: []
            },
            lastScanned: '',
            scanId: '',
            scanning: false
        };
    }

    componentDidMount() {
        this.pollBenchMarks();
    }

    componentWillUnmount() {
        if (this.timeout) this.timeout.cancel();
    }

    onTriggerScan = () => {
        this.setState({ scanning: true });
        axios.post('/v1/benchmarks/trigger', {}).then(() => {
        }).catch((error) => {
            console.log(error);
        });
    }

    onRowClick = (row) => {
        emitter.emit('ComplianceTable:row-selected', row);
    }

    getBenchMarks() {
        const promise = new Promise((resolve) => {
            const params = `?${queryString.stringify(this.params)}`;
            const { table } = this.state;
            axios.get(`/v1/benchmarks/results/grouped/cis${params}`).then((response) => {
                if (!response.data || !response.data.benchmarks) return;
                const lastScanned = dateFns.format(response.data.benchmarks[0].time, 'MM/DD/YYYY h:MM:ss A');
                const { scanId } = response.data.benchmarks[0];
                table.rows = response.data.benchmarks[0].checkResults;
                if (lastScanned !== this.state.lastScanned) {
                    this.setState({
                        table, lastScanned, scanId, scanning: false
                    });
                }
                resolve({});
            }).catch(() => {
                table.rows = [];
                this.setState({ table, lastScanned: '', scanId: '' });
                resolve({});
            });
        });
        return promise;
    }

    pollBenchMarks = () => {
        this.getBenchMarks().then(() => {
            this.timeout = setTimeout(this.pollBenchMarks, 5000);
        });
    }

    render() {
        return (
            <div className="flex flex-1">
                <Tabs className="bg-white" headers={this.state.tab.headers}>
                    <TabContent name={this.state.tab.headers[0]}>
                        <div className="flex w-full mb-3 px-3 items-center">
                            <span className="flex flex-1 text-xl font-500 text-primary-500 self-end">Last Scanned: {this.state.lastScanned || 'Never'}</span>
                            <div className="flex self-center justify-end pr-5 border-r border-primary-200">
                                <Select options={this.state.category.options} />
                            </div>
                            {
                                (this.state.scanning) ? (<button className="p-3 ml-5 h-10 w-24 rounded-sm bg-success-500 text-white hover:bg-success-600 uppercase text-center"><ClipLoader color="white" loading={this.state.scanning} size={20} /></button>) : (<button className="p-3 ml-5 h-10 w-24 rounded-sm bg-success-500 text-white hover:bg-success-600 uppercase" onClick={this.onTriggerScan}>Scan now</button>)
                            }
                        </div>
                        <div className="flex flex-1 border-t border-primary-300 bg-base-100">
                            <div className="w-full p-3 overflow-y-scroll bg-white rounded-sm shadow">

                                <Table columns={this.state.table.columns} rows={this.state.table.rows} onRowClick={this.onRowClick} />
                            </div>
                            <ComplianceBenchmarksSidePanel />
                        </div>
                    </TabContent>
                    <TabContent name={this.state.tab.headers[1]} />
                    <TabContent name={this.state.tab.headers[2]} />
                </Tabs>
            </div>
        );
    }
}

export default CompliancePage;
