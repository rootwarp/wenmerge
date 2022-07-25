import './App.css';

import BlockInfo from './components/blockInfo.js';

function App() {
  return (
    <div className="app">
      <div className="app__title_container">
        <p className='app__title'>wenmerge?</p>
        <p className='app__from'>by DSRV</p>
        <a href='https://dsrvlabs.com' className='app__home_link'>> GO HOME</a>
      </div>
      <BlockInfo title="MSF10 - Mainnet shasow fork" target_ttd="54892060024013871644672" rpc="https://api-wenmerge.dsrvlabs.dev"></BlockInfo>
      <BlockInfo title="goerli merge" target_ttd="10790000" rpc="https://api-wenmerge-goerli.dsrvlabs.dev"></BlockInfo>
    </div>
  );
}

export default App;
