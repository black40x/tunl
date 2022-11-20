import React from 'react';
import Box from '@mui/joy/Box';

function Root(props) {
    return (
        <Box
            {...props}
            sx={[
                {
                    display: 'grid',
                    gridTemplateColumns: {
                        xs: '1fr',
                        md: 'minmax(200px, 500px) minmax(500px, 1fr)',
                    },
                    gridTemplateRows: '64px 1fr',
                    height: '100vh',
                    overflow: 'hidden'
                },
                ...(Array.isArray(props.sx) ? props.sx : [props.sx]),
            ]}
        />
    );
}

function Header(props) {
    return (
        <Box
            component="header"
            className="Header"
            {...props}
            sx={[
                {
                    p: 2,
                    gap: 2,
                    bgcolor: 'background.surface',
                    display: 'flex',
                    flexDirection: 'row',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    gridColumn: '1 / -1',
                    borderBottom: '1px solid',
                    borderColor: 'divider',
                    position: 'sticky',
                    top: 0,
                    zIndex: 1100,
                },
                ...(Array.isArray(props.sx) ? props.sx : [props.sx]),
            ]}
        />
    );
}

function SidePane(props) {
    return (
        <Box
            className="Requests"
            {...props}
            sx={[
                {
                    bgcolor: 'background.surface',
                    borderRight: '1px solid',
                    borderColor: 'divider',
                    display: {
                        xs: 'none',
                        md: 'initial',
                    },
                    height: '100%',
                    overflow: 'hidden',
                    overflowY: 'auto'
                },
                ...(Array.isArray(props.sx) ? props.sx : [props.sx]),
            ]}
        />
    );
}

function Main(props) {
    return (
        <Box
            component="main"
            className="Main"
            {...props}
            sx={[{p: 2, bgcolor: 'background.body', overflowX: 'auto'}, ...(Array.isArray(props.sx) ? props.sx : [props.sx])]}
        />
    );
}

export default {
    Root,
    Header,
    SidePane,
    Main,
};