import { initializeApp } from "firebase/app";
import { getAnalytics, logEvent } from "firebase/analytics";
import './App.css';

import BlockInfo from './components/blockInfo.js';

// Import the functions you need from the SDKs you need
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: "AIzaSyA9zfAZvgJ4OsuZaEe7LCOS9hX7McZm0bY",
  authDomain: "validator-k8s.firebaseapp.com",
  projectId: "validator-k8s",
  storageBucket: "validator-k8s.appspot.com",
  messagingSenderId: "1020764648167",
  appId: "1:1020764648167:web:1ce39986afcd0f8e182842",
  measurementId: "G-3Y52RWFPEZ"
};

const app = initializeApp(firebaseConfig);
const analytics = getAnalytics(app);

function App() {
  logEvent(analytics, 'open App');
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
