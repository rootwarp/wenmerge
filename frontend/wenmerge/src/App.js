import './App.css';

import BlockInfo from './components/blockInfo.js';
import TBD from './components/tbd.js';

function App() {
  return (
    <div className="App">
      <h className='App__title'>wenmerge?</h>
      <BlockInfo title="GSF5 - goerli shadow fork" target_ttd="10625003" rpc="https://api-wenmerge-goerli.dsrvlabs.dev"></BlockInfo>
      <BlockInfo title="goerli merge" target_ttd="10790000" rpc="https://api-wenmerge-goerli.dsrvlabs.dev"></BlockInfo>
      <TBD title="mainnet"></TBD>
    </div>
  );
}

export default App;
