import { React } from 'react';

import './tbd.css';

const TBD = (props) => {

    const { title } = props;

    return (
        <div className="tbd_container">
            <div className="tbd__title">{ title }</div>
            <div className="tbd__msg">TBD</div>
        </div>
    );
}

export default TBD;