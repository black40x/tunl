import * as React from 'react';
import Box from '@mui/joy/Box';
import Typography from '@mui/joy/Typography';
import Chip from '@mui/joy/Chip';
import List from '@mui/joy/List';
import ListDivider from '@mui/joy/ListDivider';
import ListItem from '@mui/joy/ListItem';
import ListItemButton from '@mui/joy/ListItemButton';
import getMethodColor from "../helper";

export default function RequestList(props) {
    const {requests, selected, onSelect} = props;

    return (
        <List sx={{pt: 0}}>
            {requests.map((item, index) => (
                <React.Fragment key={index}>
                    <ListItem onClick={() => onSelect(item)}>
                        <ListItemButton
                            {...(item.uuid === selected && {variant: 'soft', color: 'primary'})}
                            sx={{p: 2}}
                        >
                            <Box sx={{width: '100%', userSelect: 'none', display: 'flex', alignItems: 'center'}}>
                                <Box
                                    sx={{
                                        display: 'grid',
                                        gridTemplateColumns: {
                                            xs: '90px 1fr',
                                        },
                                        justifyContent: 'space-between',
                                        alignItems: 'center'
                                    }}
                                >
                                    <Box
                                        sx={{pr: 2, textAlign: 'right'}}
                                    >
                                        <Chip
                                            size="sm"
                                            variant="outlined"
                                            sx={{textTransform: 'uppercase', "--Chip-radius": "8px"}}
                                            color={getMethodColor(item.method)}
                                        >{item.method}</Chip>
                                    </Box>
                                    <Typography
                                        level="body2"
                                        sx={{overflowWrap: 'anywhere'}}
                                    >
                                        {item.uri}
                                    </Typography>
                                </Box>
                            </Box>
                        </ListItemButton>
                    </ListItem>
                    <ListDivider sx={{m: 0}}/>
                </React.Fragment>
            ))}
        </List>
    );
}