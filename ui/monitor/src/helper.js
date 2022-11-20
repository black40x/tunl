export default function getMethodColor(method) {
    switch (method.toLowerCase()) {
        case 'post': return 'warning'
        case 'get': return 'success'
        case 'patch': return 'danger'
        default: return 'neutral'
    }
}