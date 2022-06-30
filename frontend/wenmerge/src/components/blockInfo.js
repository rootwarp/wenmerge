import { React, useState } from 'react';

import { getDifficalty } from '../lib/difficulty';

const TARGET = '52931539986079836345517';

const BlockInfo = () => {
    useState(() => {
        console.log('useState');
        getDifficalty(TARGET);
    });

    return (
        <div>
            Block Info
        </div>
    );
}

export default BlockInfo;