import { React, useEffect, useState } from 'react';

import { getDifficalty } from '../lib/difficulty';
import "./blockInfo.css"

const MAX_BLOCK_ENTRIES = 5;

const numberWithCommas = (x) => {
    if(x == null) {
        return "";
    } else {
        return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    }
}

const BlockInfo = (props) => {
    const [ ttd, setTTDState ] = useState({
        target_difficulty: null,
        difficulty_velocity: null,
        current_block_number: null,
        expect_ttd_block_number: null,
        expect_ttd_time: null,
    });

    const [ttdStates, setTTDStates] = useState([]);

    const { title, target_ttd, rpc } = props;

    useEffect(() => {
        console.log('useEffect');
        const timer = setInterval(() => {
            getDifficalty(rpc, target_ttd, d => {
                setTTDState({
                    target_difficulty: d.target_difficulty,
                    difficulty_velocity: d.difficulty_velocity,
                    current_block_number: d.current_block_number,
                    current_block_timestamp: d.current_block_timestamp,
                    expect_ttd_block_number: d.expect_ttd_block_number,
                    expect_ttd_time: d.expect_ttd_time,
                    average_block_interval: d.average_block_interval,
                    current_difficulty: d.current_difficulty,
                });

                const lastItem = ttdStates[ttdStates.length-1];
                if(ttdStates.length > 0 && lastItem.current_block_number === d.current_block_number) {
                    return;
                }
                ttdStates.push(d);
                if (ttdStates.length > MAX_BLOCK_ENTRIES) {
                    ttdStates.shift();
                }

                setTTDStates(ttdStates);
            });
        }, 1000);

        return () => {
            console.log("unmounted");
            clearInterval(timer);
        }
    }, []);

    return (
        <div className='wenmerge_container'>
            <div className='title_container'>
                <div className='title__name'>{ title == null ? "--" : title }</div>
            </div>

            <div className='summary_container'>
                <div className='summary__row'>
                    <div>Now</div>
                    <div> { new Date().toISOString() }</div>
                </div>

                <div className='summary__row'>
                    <div>Expect TTD time</div>
                    <div>{ ttd.expect_ttd_time }</div>
                </div>

                <div className='summary__row'>
                    <div>Block No.</div>
                    <div>{ numberWithCommas(ttd.current_block_number) }</div>
                </div>

                <div className='summary__row'>
                    <div>Diff. Velocity</div>
                    <div>{ numberWithCommas(ttd.difficulty_velocity) }</div>
                </div>

                <div className='summary__row'>
                    <div>Block Time</div>
                    <div>{ ttd.average_block_interval}</div>
                </div>

                <div className='summary__row'>
                    <div>Target Difficulty</div>
                    <div>{ ttd.target_difficulty == null ? '' : ttd.target_difficulty.toLocaleString('fullwide') }</div>
                </div>

                <div className='summary__row'>
                    <div>Expect Block No</div>
                    <div>{ numberWithCommas(ttd.expect_ttd_block_number) }</div>
                </div>
            </div>

            <div className='blockinfo_container'>
                <div className='blockinfo__title'>
                    <div>Block No.</div>
                    <div>Time Created</div>
                    <div>Difficulty</div>
                </div>

                <div className='blockinfo_list'>
                    {
                        ttdStates.map(item => {
                            return(
                                <div key={ item.current_block_number } className='blockinfo__detail'>
                                    <div>{ numberWithCommas(item.current_block_number) }</div>
                                    <div>{ item.current_block_timestamp }</div>
                                    <div>{ numberWithCommas(item.current_difficulty) }</div>
                                </div>
                            );
                        }).reverse()
                    }
                </div>
            </div>
        </div>
    );
}

export default BlockInfo;