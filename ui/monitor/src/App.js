import './App.scss';
import {useEffect, useState} from 'react'
import RequestPanel from "./components/RequestPanel";
import RequestList from "./components/RequestList";
import RequestDetail from "./components/RequestDetail";

function App() {
    const localPort = 6060
    const [requests, setRequests] = useState([])
    const [filterStr, setFilterStr] = useState('')
    const [currentRequest, setCurrentRequest] = useState(undefined)
    const [connectStatus, setConnectStatus] = useState(false)
    let connectInterval = undefined

    const connectToLocal = () => {
        let ws = new WebSocket(`ws://localhost:${localPort}/connect`)
        ws.onopen = function(evt) {
            setConnectStatus(true)

            if (connectInterval !== undefined) {
                clearInterval(connectInterval)
            }
        }
        ws.onclose = function(evt) {
            ws = null
            setConnectStatus(false)
            tryToConnect()
        }
        ws.onmessage = function(evt) {
            setRequests(o => [JSON.parse(evt.data), ...o])
        }
    }

    const tryToConnect = () => {
        if (!connectInterval) {
            connectInterval = setInterval(connectToLocal, 1000)
        }
    }

    const selectRequest = (r) => {
        setCurrentRequest(r)
    }

    const filter = (f) => {
        setFilterStr(f)
    }

    const filterRequests = () => {
        if (filterStr === '') {
            return requests
        }
        return requests.filter(i => i.uri.includes(filterStr))
    }

    useEffect(() => {
        connectToLocal()
    }, [])

    return (
        <div className="App">
            <RequestPanel onFilter={filter} connected={connectStatus}>
                <RequestList
                    requests={filterRequests()}
                    selected={currentRequest?.uuid}
                    onSelect={selectRequest} />
            </RequestPanel>
            { currentRequest !== undefined && <RequestDetail request={currentRequest} /> }
        </div>
    );
}

export default App;
