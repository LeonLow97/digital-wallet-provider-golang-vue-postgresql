const formatDateTime = (dateTime) => {
  const options = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }
  return new Date(dateTime).toLocaleString('en-US', options)
}

export default formatDateTime
