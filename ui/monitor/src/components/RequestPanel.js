import './RequestPanel.scss'

function RequestPanel(props) {
    const { children, connected, onFilter } = props

    const filter = (e) => {
        if (e.keyCode === 13 && onFilter) {
            onFilter(e.target.value)
        }
    }

    return (
        <div className={'RequestPanel'}>
            <div className={'RequestPanel__Header'}>
                <span className={'RequestPanel__Logo'}>
                    tunl monitor
                    <span className={'RequestPanel__Status ' + (
                        connected ? 'RequestPanel__Status--Online' : 'RequestPanel__Status--Offline'
                    )}>{connected ? 'online' : 'offline'}</span>
                </span>
                <input
                    className={'RequestPanel__Filter'}
                    type="text"
                    placeholder="Filter and press Enter"
                    onKeyDown={filter}/>
            </div>
            {children}
        </div>
    )
}

export default RequestPanel