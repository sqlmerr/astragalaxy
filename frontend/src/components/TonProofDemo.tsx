import { Button, Typography } from 'antd';
import {useCallback, useEffect, useRef, useState} from 'react';
import ReactJson from 'react-json-view';
import { TonProofDemoApi } from '../TonProofApi';
import {useTonConnectUI, useTonWallet} from "@tonconnect/ui-react";
import useInterval from "../hooks/useInterval.ts";

const { Title } = Typography;

export function TonProofDemo() {
    const firstProofLoading = useRef<boolean>(true);

    const [data, setData] = useState({});
    const wallet = useTonWallet();
    const [authorized, setAuthorized] = useState(false);
    const [tonConnectUI] = useTonConnectUI();

    const handleClick = useCallback(async () => {
        if (!wallet) {
            return;
        }
        const response = await TonProofDemoApi.getAccountInfo(wallet.account);

        setData(response);
    }, [wallet]);

    const recreateProofPayload = useCallback(async () => {
        if (firstProofLoading.current) {
            tonConnectUI.setConnectRequestParameters({ state: 'loading' });
            firstProofLoading.current = false;
        }

        const payload = await TonProofDemoApi.generatePayload();
        console.log(payload)

        if (payload) {
            tonConnectUI.setConnectRequestParameters({ state: 'ready', value: {tonProof: payload} });
        } else {
            tonConnectUI.setConnectRequestParameters(null);
        }
    }, [tonConnectUI, firstProofLoading])

    if (firstProofLoading.current) {
        recreateProofPayload();
    }

    useInterval(recreateProofPayload, 9 * 60 * 1000);

    useEffect(() =>
        tonConnectUI.onStatusChange(async w => {
            if (!w) {
                TonProofDemoApi.reset();
                setAuthorized(false);
                return;
            }

            if (w.connectItems?.tonProof && 'proof' in w.connectItems.tonProof) {
                await TonProofDemoApi.checkProof(w.connectItems.tonProof.proof, w.account);
            }

            if (!TonProofDemoApi.accessToken) {
                await tonConnectUI.disconnect();
                setAuthorized(false);
                return;
            }

            setAuthorized(true);
        }), [tonConnectUI]);

    if (!authorized) {
        return <h1>unauthorized</h1>
    }
    console.log(wallet);

    return (
        <div className="ton-proof-demo">
        <Title level={3}>Demo backend API with ton_proof verification</Title>
    {wallet ? (
        <Button type="primary" shape="round" onClick={handleClick}>
        Call backend getAccountInfo()
    </Button>
    ) : (
        <div className="ton-proof-demo__error">Connect wallet to call API</div>
    )}
    <ReactJson src={data} name="response" theme="ocean" />
        </div>
);
}