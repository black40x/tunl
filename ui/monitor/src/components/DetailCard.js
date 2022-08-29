import './DetailCard.scss'

function DetailCard(props) {
    const { label, children } = props

    return (
        <div className={'DetailCard'}>
            <div className={'DetailCard__Header'}>{label}</div>
            <div className={'DetailCard__Body'}>
                {children}
            </div>
        </div>
    )
}

export default DetailCard