import './App.css';

import Banner from './components/Banner/Banner';
import ChatWidget from './components/ChatWidget/ChatWidget';
import PageFooter from './components/Footer/Footer';
import Main from './components/Main/Main';

const App = () => {
  return (
    <div className="App">
      <Banner />
      <Main />
      <PageFooter />
      <ChatWidget />
    </div>
  );
};

export default App;
