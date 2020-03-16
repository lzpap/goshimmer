import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {Provider} from 'mobx-react';
import {createBrowserHistory} from 'history';
import 'chartjs-plugin-streaming';
import {App} from 'app/App';
import {RouterStore, syncHistoryWithStore} from 'mobx-react-router';
import {Router} from 'react-router-dom';
import NodeStore from "app/stores/NodeStore";
import ExplorerStore from "app/stores/ExplorerStore";
import FaucetStore from "app/stores/FaucetStore";

// prepare MobX stores
const routerStore = new RouterStore();
const nodeStore = new NodeStore();
const explorerStore = new ExplorerStore(routerStore);
const faucetStore = new FaucetStore();
const stores = {
    "routerStore": routerStore,
    "nodeStore": nodeStore,
    "explorerStore": explorerStore,
    "faucetStore": faucetStore,
};

const browserHistory = createBrowserHistory();
const history = syncHistoryWithStore(browserHistory, routerStore);

// render react DOM
ReactDOM.render(
    <Provider {...stores}>
        <Router history={history}>
            <App/>
        </Router>
    </Provider>,
    document.getElementById('root')
);
