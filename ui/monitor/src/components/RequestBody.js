import './RequestBody.scss'

function RequestBody(props) {
    const {body, bodyType} = props

    console.log(body)
    console.log(bodyType)

    if (!bodyType || !body) {
        return (
            <pre className={'RequestBody'}>Body empty</pre>
        )
    } else if (bodyType === 'json') {
        return (
            <pre className={'RequestBody'}>{body}</pre>
        )
    } else if (bodyType === 'form-data') {
        return (
            <table className={'RequestBodyTable'}>
                <thead>
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                </tr>
                </thead>
                <tbody>
                {body.values && Object.keys(body.values).map((k, v) => {
                    return (<tr key={k}><td>{k}</td><td>{body.values[k]}</td></tr>)
                })}
                {body.files && Object.keys(body.files).map((k, v) => {
                    return (<tr key={k}><td>{k}</td><td>{body.files[k]}</td></tr>)
                })}
                </tbody>
            </table>
        )
    }
}

export default RequestBody