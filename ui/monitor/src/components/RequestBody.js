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
            <Typography level="body2" component="pre">{ body }</Typography>
        )
    } else if (bodyType === 'form-data') {
        return (
            <Box sx={{
                mb: .3,
                display: 'grid',
                gridTemplateColumns: '1fr 1fr',
                gap: 1
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
                {body.values && Object.keys(body.files).map((k, v) => {
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