import React from "./_snowpack/pkg/react.js";
import ReactDOM from "./_snowpack/pkg/react-dom.js";
import Autocomplete, {createFilterOptions} from "./_snowpack/pkg/@mui/material/Autocomplete.js";
import TextField from "./_snowpack/pkg/@mui/material/TextField.js";
import {createChart} from "./_snowpack/pkg/lightweight-charts.js";
class StockChart extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      options: [],
      symbol: this.props.symbol
    };
    this.chartContainer = null;
    this.candlestickSeries = null;
    this.pickerFilterOptions = createFilterOptions({matchFrom: "start"});
    this.setChartContainerRef = (element) => this.chartContainer = element;
    this.onPickerChange = (_, newSymbol) => {
      this.setState({symbol: newSymbol});
      this.populateChartSeries(newSymbol);
    };
    fetch("/symbols").then((response) => response.json()).then((symbols) => {
      this.setState({loading: false, options: symbols});
    });
  }
  componentDidMount() {
    const chart = createChart(this.chartContainer);
    this.candlestickSeries = chart.addCandlestickSeries();
    if (this.state.symbol) {
      this.populateChartSeries(this.state.symbol);
    }
  }
  populateChartSeries(symbol) {
    fetch(`/prices/${symbol}`).then((response) => response.json()).then((series) => this.candlestickSeries.setData(series));
  }
  render() {
    return /* @__PURE__ */ React.createElement("div", {
      id: "stock-chart"
    }, /* @__PURE__ */ React.createElement(Autocomplete, {
      disableClearable: true,
      selectOnFocus: true,
      filterOptions: this.pickerFilterOptions,
      loading: this.state.loading,
      onChange: this.onPickerChange,
      options: this.state.options,
      sx: {width: 110},
      value: this.state.symbol,
      renderInput: (params) => /* @__PURE__ */ React.createElement(TextField, {
        ...params,
        label: "Symbol"
      })
    }), /* @__PURE__ */ React.createElement("div", {
      id: "lightweight-chart-container",
      ref: this.setChartContainerRef
    }));
  }
}
ReactDOM.render(/* @__PURE__ */ React.createElement(StockChart, null), document.getElementById("root"));
