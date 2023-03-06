import './App.css';

import Banner from './components/Banner/Banner';
import PageFooter from './components/Footer/Footer';
import Main from './components/Main/Main';

const App = () => {
  return (
    <div className="App">
      <Banner />
      <Main />
      <PageFooter />
    </div>
  );
};

export default App;
