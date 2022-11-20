import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import { StyledEngineProvider } from '@mui/joy/styles';
import "./index.css"

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <StyledEngineProvider injectFirst>
        <App />
    </StyledEngineProvider>
);
