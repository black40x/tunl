import RequestListItem from './RequestListItem'
import './RequestList.scss'

function RequestLst(props) {
    const { requests, selected, onSelect } = props

    return (
        <div className={'RequestList'}>
            {requests.map((r, k) => (
                <RequestListItem selected={selected === r.uuid} key={k} request={r} onSelect={onSelect} />
            ))}
        </div>
    )
}

export default RequestLst