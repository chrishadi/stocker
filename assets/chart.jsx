/* Add JavaScript code here! */
import React from 'react'
import ReactDOM from 'react-dom'
import Autocomplete, { createFilterOptions } from '@mui/material/Autocomplete'
import TextField from '@mui/material/TextField'
import { createChart } from 'lightweight-charts'
import FibonacciPivot from './fibonacci_pivot'

class StockChart extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      loading: true,
      lastDateData: null,
      options: [],
      symbol: this.props.symbol,
    }

    this.chartContainer = null
    this.candlestickSeries = null
    this.pickerFilterOptions = createFilterOptions({matchFrom: 'start'})

    this.setChartContainerRef = element => this.chartContainer = element
    this.onPickerChange = (_, newSymbol) => {
      this.setState({loading: true})
      this.updateChartSeries(newSymbol)
    }

    fetch('/symbols')
    .then(response => response.json())
    .then(symbols => { this.setState({ loading: false, options: symbols }) } )
  }

  componentDidMount() {
    const chart = createChart(this.chartContainer)
    this.candlestickSeries = chart.addCandlestickSeries()
    if (this.state.symbol) {
      this.updateChartSeries(this.state.symbol)
    }
  }

  updateChartSeries(symbol) {
    fetch(`/prices/${symbol}`)
    .then(response => response.json())
    .then(series => {
      this.setState({ loading: false, symbol: symbol, lastDateData: series[series.length-1] })
      this.candlestickSeries.setData(series)
    })
  }

  render() {
    return (
     <div id="stock-chart">
        <Autocomplete
          disableClearable
          selectOnFocus
          filterOptions={this.pickerFilterOptions}
          loading={this.state.loading}
          onChange={this.onPickerChange}
          options={this.state.options}
          sx={{ width: 110 }}
          value={this.state.symbol}
          renderInput={(params) => <TextField {...params} label="Symbol" />}
        />
        <div id="lightweight-chart-container" ref={this.setChartContainerRef} />
        <FibonacciPivot data={this.state.lastDateData} />
      </div>
    )
  }
}

ReactDOM.render(<StockChart />, document.getElementById('root'))
