import './RequestListItem.scss'
import RequestMethod from './RequestMethod'

function RequestListItem(props) {
    const { request, selected, onSelect } = props

    function onClick () {
        if (onSelect) {
            onSelect(request)
        }
    }

    return (
        <div className={'RequestListItem' + (selected ? ' RequestListItem--Selected' : '')} onClick={onClick}>
            <span className={'RequestListItem__Method'}><RequestMethod method={request.method}/></span>
            <span className={'RequestListItem__Uri'}>{request.uri}</span>
        </div>
    )
}

export default RequestListItem