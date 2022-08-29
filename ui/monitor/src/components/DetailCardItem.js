import './DetailCardItem.scss'

function DetailCardItem(props) {
    const { label, value } = props
    return (
        <div className={'DetailCardItem'}>
            <span className={'DetailCardItem__Label'}>{label}</span>
            <span className={'DetailCardItem__Value'}>{value}</span>
        </div>
    )
}

export default DetailCardItem