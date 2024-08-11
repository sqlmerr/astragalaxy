import {THEME, TonConnectUIProvider} from "@tonconnect/ui-react";
import Header from "./components/Header";
import { TonProofDemo } from "./components/TonProofDemo";

export default function App() {
    return (
        <TonConnectUIProvider
            manifestUrl="https://astragalaxy.vercel.app/tonconnect-manifest.json"
            uiPreferences={{
                theme: THEME.DARK,
            }}

        >
            <div className="app">
                    <Header/>
                    <TonProofDemo/>
            </div>
        </TonConnectUIProvider>
    );
}
