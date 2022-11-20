import * as React from 'react';
import Box from '@mui/joy/Box';
import Typography from '@mui/joy/Typography';

function RequestBody(props) {
    const {body, bodyType} = props

    if (!bodyType || !body) {
        return (
            <Typography level="body2" component="pre">Body empty</Typography>
        )
    } else if (bodyType === 'json') {
        return (
            <Typography
                sx={{
                    borderRadius: 'sm',
                    p: 2,
                    display: 'block',
                    background: 'var(--joy-palette-neutral-800, #001E3C)',
                    borderColor: 'var(--joy-palette-neutral-700, #132F4C)',
                    borderWidth: '1px',
                    borderStyle: 'solid',
                    fontFamily: 'Menlo,Consolas,"Droid Sans Mono",monospace',
                    fontSize: '0.8125rem'
                }}
                level="body2"
                component="pre">
                { JSON.stringify(JSON.parse(body), null, 2) }
            </Typography>
        )
    }  else if (bodyType === 'xml' || bodyType === 'text') {
        return (
            <Typography
                sx={{
                    borderRadius: 'sm',
                    p: 2,
                    display: 'block',
                    background: 'var(--joy-palette-neutral-800, #001E3C)',
                    borderColor: 'var(--joy-palette-neutral-700, #132F4C)',
                    borderWidth: '1px',
                    borderStyle: 'solid',
                    fontFamily: 'Menlo,Consolas,"Droid Sans Mono",monospace',
                    fontSize: '0.8125rem'
                }}
                level="body2"
                component="pre">
                { body }
            </Typography>
        )
    } else if (bodyType === 'form-data' || bodyType === 'url-encoded') {
        return (
            <Box sx={{
                mb: .3,
                display: 'grid',
                gridTemplateColumns: '1fr 1fr',
                gap: 1,
                borderRadius: 'sm',
                p: 2,
                background: 'var(--joy-palette-neutral-800, #001E3C)',
                borderColor: 'var(--joy-palette-neutral-700, #132F4C)',
                borderWidth: '1px',
                borderStyle: 'solid',
                fontFamily: 'Menlo,Consolas,"Droid Sans Mono",monospace',
                fontSize: '0.8125rem'
            }}>
                <Typography level="body3">KEY</Typography>
                <Typography level="body3">VALUE</Typography>
                {body.values && Object.keys(body.values).map((k, v) => {
                    return (
                        <React.Fragment key={k}>
                            <Typography level="body2">{k}</Typography>
                            <Typography level="body2">{ body.values[k] }</Typography>
                        </React.Fragment>
                    )
                })}
                {body.files && Object.keys(body.files).map((k, v) => {
                    return (
                        <React.Fragment key={k}>
                            <Typography level="body2">{k}</Typography>
                            <Typography level="body2">{ body.files[k] }</Typography>
                        </React.Fragment>
                    )
                })}
            </Box>
        )
    }
}

export default RequestBody