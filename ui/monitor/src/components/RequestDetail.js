import './RequestDetail.scss'
import DetailCard from './DetailCard'
import DetailCardItem from './DetailCardItem'
import React from 'react'
import RequestMethod from './RequestMethod'
import RequestBody from './RequestBody'

function RequestDetail(props) {
    const { request } = props

    const headers = () => {
        let hArr = []
        request.header.forEach(head => {
            head.value.forEach(val => {
                hArr.push({
                    key: head.key,
                    value: val
                })
            })
        })

        return hArr
    }

    return (
        <div className={'RequestDetail'}>
            <DetailCard label={'General'}>
                <DetailCardItem label={'Request URL'} value={request.uri} />
                <DetailCardItem label={'Request Method'} value={
                    <React.Fragment><RequestMethod method={request.method} /></React.Fragment>
                } />
                <DetailCardItem label={'Request Proto'} value={request.proto} />
                <DetailCardItem label={'Remote Address'} value={request.remote_address} />
                <DetailCardItem label={'Duration'} value={`${request.duration}ms`} />
            </DetailCard>

            <DetailCard label={'Request Headers'}>
                {headers().map((h, k) => (
                    <DetailCardItem key={k} label={h.key} value={h.value} />
                ))}
            </DetailCard>

            {request.body && (
                <DetailCard label={'Request Body'}>
                    <RequestBody body={request.body} bodyType={request.body_type} />
                </DetailCard>
            )}
        </div>
    )
}

export default RequestDetail