import './App.css';

import BlockInfo from './components/blockInfo.js';
import TBD from './components/tbd.js';

function App() {
  return (
    <div className="App">
      <h className='App__title'>wenmerge?</h>
      <TBD title="mainnet"></TBD>
      <BlockInfo title="goerli shadow fork" target_ttd="10627192" rpc="https://api-wenmerge-goerli.dsrvlabs.dev"></BlockInfo>
      <BlockInfo title="goerli merge" target_ttd="10790000" rpc="https://api-wenmerge-goerli.dsrvlabs.dev"></BlockInfo>
    </div>
  );
}

export default App;
