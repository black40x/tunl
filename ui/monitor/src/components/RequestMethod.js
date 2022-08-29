import './RequestMethod.scss'

function RequestMethod(props) {
    const { method } = props
    const methodClass = () => {
        let className = 'RequestMethod'
        if (method.toLowerCase() === 'post') {
            className += ' RequestMethod--Post'
        }
        if (method.toLowerCase() === 'get') {
            className += ' RequestMethod--Get'
        }

        return className
    }

    return (
        <span className={methodClass()}>{method}</span>
    )
}

export default RequestMethod