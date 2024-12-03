import ForceGraph from './components/SitesGraph/SitesGraph';
import {GraphData} from "./types/GraphData.ts";

function App() {
    const bib: GraphData = {
        nodes: [
            {id: 'site1'},
            {id: 'site2'},
            {id: 'site3'},
            {id: 'site4'},
            {id: 'site5'},
            {id: 'site6'},
        ],
        links: [
            {source: 'site1', target: 'site2'},
            {source: 'site1', target: 'site3'},
            {source: 'site2', target: 'site3'},
            {source: 'site4', target: 'site5'},
            {source: 'site4', target: 'site6'},
            {source: 'site5', target: 'site6'},
            {source: 'site1', target: 'site6'},
        ],
    }

    return (
        <div className="App">
            <div>
                <ForceGraph width={1000} data={bib} backgroundCol={'#020202'} height={1000}/>
            </div>
        </div>
    );
}

export default App;
