import React from 'react'
import PropTypes from 'prop-types'

export default function FibonacciPivot(props) {
  var pp, r1, r2, r3, s1, s2, s3

  if (props.data) {
    var d = props.data.high - props.data.low
    pp = (props.data.high + props.data.low + props.data.close) / 3
    r1 = (pp + 0.382*d).toFixed(2)
    r2 = (pp + 0.618*d).toFixed(2)
    r3 = (pp + d).toFixed(2)
    s1 = (pp - 0.382*d).toFixed(2)
    s2 = (pp - 0.618*d).toFixed(2)
    s3 = (pp - d).toFixed(2)
    pp = pp.toFixed(2)
  }

  return (
    <div id="fibonacci-pivot">
      <span>Fibonacci Pivot:</span>&nbsp;
      <span>{ s3 }</span>&nbsp;
      <span>{ s2 }</span>&nbsp;
      <span>{ s1 }</span>&nbsp;
      <span>{ pp }</span>&nbsp;
      <span>{ r1 }</span>&nbsp;
      <span>{ r2 }</span>&nbsp;
      <span>{ r3 }</span>
    </div>
  )
}

FibonacciPivot.PropTypes = {
    data: PropTypes.object
}

FibonacciPivot.defaultProps = {
    data: null
}
