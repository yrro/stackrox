import React from 'react';
import Tabs from 'Components/Tabs';
import TabContent from 'Components/TabContent';

import BenchmarksPage from 'Containers/Compliance/BenchmarksPage';

const CompliancePage = () => {
    const benchmarks = [
        { text: 'CIS Benchmark', disabled: false },
        { text: 'Swarm Benchmark', disabled: false },
        { text: 'Kubernetes v1.2.0 Benchmark', disabled: false }
    ];
    return (
        <section className="flex flex-1 h-full">
            <div className="flex flex-1">
                <Tabs className="bg-white" headers={benchmarks}>
                    {
                        benchmarks.map(benchmark => (
                            <TabContent key={benchmark.text}>
                                <BenchmarksPage benchmarkName={benchmark.text} />
                            </TabContent>
                        ))
                    }
                </Tabs>
            </div>
        </section>
    );
};

export default CompliancePage;
