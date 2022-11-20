import * as React from 'react';
import Box from '@mui/joy/Box';
import Sheet from '@mui/joy/Sheet';
import Typography from '@mui/joy/Typography';
import Divider from '@mui/joy/Divider';
import RequestBody from "./RequestBody";
import getMethodColor from "../helper";
import Chip from "@mui/joy/Chip";

export default function RequestContent(props) {
    const { request } = props;

    const headers = () => {
        let hArr = [];
        request?.header.forEach(head => {
            head.value.forEach(val => {
                hArr.push({
                    key: head.key,
                    value: val
                });
            });
        });
        return hArr;
    }

    return (
        <Sheet
            variant="outlined"
            sx={{
                minHeight: 500,
                borderRadius: 'sm',
                p: 2,
                mb: 3,
            }}
        >
            <Box
                sx={{ display: 'flex', flexDirection: 'row', alignItems: 'center', gap: 2 }}
            >
                <Chip
                    size="md"
                    variant="outlined"
                    sx={{textTransform: 'uppercase', "--Chip-radius": "8px"}}
                    color={getMethodColor(request?.method)}
                >{request?.method}</Chip>

                <Typography level="h5" textColor="text.primary">
                    { request?.uri }
                </Typography>
            </Box>
            <Divider sx={{mt: 2}} />
            <Typography component="div" level="body2" mt={2} mb={2}>
                <Typography level="body1" sx={{ display: 'block', fontWeight: 500, mb: 1 }}>General</Typography>
                <Box sx={{
                    mb: .3,
                    display: 'grid',
                    gridTemplateColumns: 'minmax(100px, 150px) 1fr',
                    gap: 1
                }}>
                    <Typography level="body2">Request URL:</Typography>
                    <Typography level="body2" sx={{ opacity: .7, overflowWrap: 'anywhere' }}>{ request?.uri }</Typography>
                    <Typography level="body2">Request Method:</Typography>
                    <Typography level="body2" sx={{ opacity: .7, overflowWrap: 'anywhere' }}>{ request?.method }</Typography>
                    <Typography level="body2">Request Proto:</Typography>
                    <Typography level="body2" sx={{ opacity: .7, overflowWrap: 'anywhere' }}>{ request?.proto }</Typography>
                    <Typography level="body2">Remote Address:</Typography>
                    <Typography level="body2" sx={{ opacity: .7, overflowWrap: 'anywhere' }}>{ request?.remote_address }</Typography>
                </Box>

                <Typography level="body1" sx={{ display: 'block', fontWeight: 500, mb: 1, mt: 2 }}>Request Headers</Typography>
                <Box sx={{
                    mb: .3,
                    display: 'grid',
                    gridTemplateColumns: 'minmax(100px, 150px) 1fr',
                    gap: 1
                }}>
                    {headers().map((h, k) => (
                        <React.Fragment key={k}>
                            <Typography level="body2" sx={{ mr: 1 }}>{ h.key }:</Typography>
                            <Typography level="body2" sx={{ opacity: .7, overflowWrap: 'anywhere' }}>{ h.value }</Typography>
                        </React.Fragment>
                    ))}
                </Box>

                {request?.body && (
                    <Box>
                        <Typography level="body1" sx={{ display: 'block', fontWeight: 500, mb: 1, mt: 2 }}>Request Body</Typography>
                        <Box sx={{ mb: .3 }}>
                            <RequestBody body={request.body} bodyType={request.body_type} />
                        </Box>
                    </Box>
                )}
            </Typography>
        </Sheet>
    );
}
