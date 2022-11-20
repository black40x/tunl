import { CssVarsProvider } from '@mui/joy/styles';
import Layout from './components/Layout';
import Box from '@mui/joy/Box';
import IconButton from '@mui/joy/IconButton';
import Typography from '@mui/joy/Typography';
import Input from '@mui/joy/Input';
import Link from '@mui/joy/Link';
import Requests from "./components/Requests";
import RequestContent from "./components/Detail";
import Badge from '@mui/joy/Badge';
import monitorTheme from './theme';
import GitHubIcon from '@mui/icons-material/GitHub';
import SearchRoundedIcon from '@mui/icons-material/SearchRounded';
import {useEffect, useState} from 'react'
import Logo from "./components/Logo";

function App() {
    const defaultPort = 6060;
    const [isOnline, setConnectStatus] = useState(false);
    const [filterStr, setFilterStr] = useState('');
    const [searchStr, setSearchStr] = useState('');
    const [currentRequest, setCurrentRequest] = useState(undefined);
    const [requests, setRequests] = useState([]);

    let connectInterval = undefined;

    const connectToLocal = () => {
        const body = document.getElementById("body")
        let localPort = parseInt(body.dataset.port)
        if (isNaN(localPort)) {
            localPort = defaultPort
        }
        let localHost = body.dataset.host

        let ws = new WebSocket(`ws://${localHost}:${localPort}/connect`)
        ws.onopen = function(evt) {
            setConnectStatus(true);

            if (connectInterval !== undefined) {
                clearInterval(connectInterval);
            }
        }
        ws.onclose = function(evt) {
            ws = null;
            setConnectStatus(false);
        }
        ws.onmessage = function(evt) {
            setRequests(o => [JSON.parse(evt.data), ...o]);
        }
    }

    const selectRequest = (r) => {
        setCurrentRequest(r);
    }

    const filterRequests = () => {
        if (filterStr === '') {
            return requests;
        }
        return requests.filter(i => i.uri.includes(filterStr));
    }

    useEffect(() => {
        connectToLocal();
    }, []);

    return (
        <CssVarsProvider theme={monitorTheme}
                         defaultMode="system"
                         modeStorageKey="identify-system-mode">
            <Layout.Root>
                <Layout.Header>
                    <Box
                        sx={{
                            display: 'flex',
                            flexDirection: 'row',
                            alignItems: 'center',
                            gap: 1.5,
                        }}
                    >
                        <Link sx={{ display: { sm: 'inline-flex' }, width: '30px' }}>
                            <Logo/>
                        </Link>
                        <Badge
                            badgeContent={isOnline ? "online" : "offline"}
                            color={isOnline ? "success" : "danger"}
                            size="sm"
                            badgeInset="0 -15px 0 0"
                        >
                            <Typography component="h1" fontWeight="xl">
                                Monitor
                            </Typography>
                        </Badge>
                    </Box>

                    <Box sx={{
                        display: 'flex',
                        flexBasis: '300px',
                        gap: 2
                    }}>
                        <Input
                            size="sm"
                            placeholder="Search request..."
                            value={searchStr}
                            onChange={(e) => {
                                setSearchStr(e.target.value)
                            }}
                            startDecorator={<SearchRoundedIcon color="primary" />}
                            endDecorator={
                                <IconButton
                                    onClick={(e) => {
                                        setSearchStr('')
                                        setFilterStr('')
                                    }}
                                    variant="outlined"
                                    size="sm"
                                    color="neutral">
                                    <Typography fontWeight="md" fontSize="sm" textColor="text.tertiary">
                                        clear
                                    </Typography>
                                </IconButton>
                            }
                            onKeyUp={(e) => {
                                if (e.code === "Enter") {
                                    setFilterStr(searchStr)
                                }
                            }}
                            sx={{
                                flexBasis: '400px',
                                display: {
                                    xs: 'none',
                                    sm: 'flex',
                                },
                            }}
                        />
                        <Link
                            variant="outlined"
                            aria-labelledby="heading-demo"
                            href="https://github.com/black40x/tunl-cli"
                            target="_blank"
                            fontSize="md"
                            borderRadius="sm"
                        >
                            <GitHubIcon />
                        </Link>
                    </Box>
                </Layout.Header>
                <Layout.SidePane>
                    <Requests
                        requests={filterRequests()}
                        selected={currentRequest?.uuid}
                        onSelect={selectRequest}
                    />
                </Layout.SidePane>
                <Layout.Main>
                    { currentRequest !== undefined && (
                        <RequestContent request={currentRequest} />
                    )}
                </Layout.Main>
            </Layout.Root>
        </CssVarsProvider>
    );
}

export default App;
