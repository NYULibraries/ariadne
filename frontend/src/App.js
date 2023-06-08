import './App.css';

import Banner from './components/Banner/Banner';
import Main from './components/Main/Main';
import PageFooter from './components/Footer/Footer';

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
