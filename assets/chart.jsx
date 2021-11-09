/* Add JavaScript code here! */
import React from 'react'
import ReactDOM from 'react-dom'
import Autocomplete, { createFilterOptions } from '@mui/material/Autocomplete'
import TextField from '@mui/material/TextField'
import { createChart } from 'lightweight-charts'

class StockChart extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      loading: true,
      options: [],
      symbol: this.props.symbol,
    }

    this.chartContainer = null
    this.candlestickSeries = null
    this.pickerFilterOptions = createFilterOptions({matchFrom: 'start'})

    this.setChartContainerRef = element => this.chartContainer = element
    this.onPickerChange = (_, newSymbol) => {
      this.setState({symbol: newSymbol})
      this.populateChartSeries(newSymbol)
    }

    fetch('/symbols')
    .then(response => response.json())
    .then(symbols => { this.setState({ loading: false, options: symbols }) } )
  }

  componentDidMount() {
    const chart = createChart(this.chartContainer)
    this.candlestickSeries = chart.addCandlestickSeries()
    if (this.state.symbol) {
      this.populateChartSeries(this.state.symbol)
    }
  }

  populateChartSeries(symbol) {
    fetch(`/prices/${symbol}`)
    .then(response => response.json())
    .then(series => this.candlestickSeries.setData(series))
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
      </div>
    )
  }
}

ReactDOM.render(<StockChart />, document.getElementById('root'))
