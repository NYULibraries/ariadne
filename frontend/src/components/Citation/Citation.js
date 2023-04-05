import "./Citation.css";

import { useEffect } from 'react';

const Citation = () => {

  const params = (new URL(window.location)).searchParams;

  // thanks to umlaut: https://github.com/NYULibraries/umlaut/blob/master/config/locales/en.yml#L28-L43
  const genresDisplayText = {
    book: "Book",
    bookitem: "Book Chapter",
    conference: "Conference",
    proceeding: "Proceeding",
    report: "Report",
    document: "Document",
    journal: "Journal",
    issue: "Issue",
    article: "Article",
    preprint: "Pre-print",
    dissertation: "Dissertation",
    unknown: ""
  }

  // prefer rft prefixed params when we have them
  // returns null if param is empty (including empty string)
  const getOpenUrlParam = (paramName) => {
    return params.get("rft." + paramName) || params.get(paramName);
  }

  const citation = {
    genre: getOpenUrlParam("genre"),
    volume: getOpenUrlParam("volume"),
    issue: getOpenUrlParam("issue"),
    start_page: getOpenUrlParam("spage"),
    end_page: getOpenUrlParam("epage"),
    pub: getOpenUrlParam("pub"),
    issn: getOpenUrlParam("issn"),
    isbn: getOpenUrlParam("isbn"),
    date: getOpenUrlParam("date"),
    doi: getOpenUrlParam("doi"),
  };

  // author either specified as "au", a series of separate params ("aufirst", "aulast", "auinit", "auinit1", auinitm"), 
  // or "aucorp"
  // thanks to umlaut: https://github.com/NYULibraries/umlaut/blob/master/app/models/referent.rb#L320-L340
  const getAuthorDisplayText = () => {
    let author;
    const aulast = getOpenUrlParam('aulast');
    const aufirst = getOpenUrlParam('aufirst');
    const auinit = getOpenUrlParam('auinit');
    const auinit1 = getOpenUrlParam('auinit1');
    const auinitm = getOpenUrlParam('auinitm');
    const aucorp = getOpenUrlParam('aucorp');
    const au = getOpenUrlParam('au');

    if (au) {
      return au;
    } else if (aulast) {
      author = aulast;
      if (aufirst) {
        author += `, ${aufirst}`;
      } else if (auinit) {
        author += `, ${auinit}`;
      } else if (auinit1) {
        author += `, ${auinit1}`;
        if (auinitm) {
          author += auinitm;
        }
      }
      return author;
    } else if (aucorp) {
      return aucorp;
    }
    return author;
  };




  citation.author = getAuthorDisplayText();

  // if we have atitle, assume we need a container title; otherwise, no container needed
  // logic from: https://github.com/NYULibraries/umlaut/blob/master/app/models/referent.rb#L288-L303
  if (getOpenUrlParam("atitle")) {
    citation.item_title = getOpenUrlParam("atitle");
    citation.container_title = getOpenUrlParam("title") || getOpenUrlParam("btitle") || getOpenUrlParam("jtitle");
  } else {
    citation.item_title = getOpenUrlParam("title") || getOpenUrlParam("btitle") || getOpenUrlParam("jtitle");
  }

  // set page title based on item title
  useEffect(() => {
    if (citation.item_title)
      document.title = 'GetIt | ' + citation.item_title;
    else
      document.title = 'GetIt';
  }, [citation.item_title]);

  const renderCitation = (citation) => {
    if (citation.container_title || citation.volume || citation.issue || citation.start_page || citation.end_page) {
      return (
        <p>
          <span className="published-in">{citation.container_title && 'Published in '}</span>
          <span className="container-title">{citation.container_title && citation.container_title + '. '}</span>
          {citation.volume && 'Volume ' + citation.volume + '. '}
          {citation.issue && 'Issue ' + citation.issue + '. '}
          {citation.start_page && 'Page ' + citation.start_page}
          {citation.end_page && '-' + citation.end_page + '. '}
        </p>
      );
    }
    return null;
  };

  return (
    <div className="citation-container">
      {citation.genre && <p className="resource-type">{genresDisplayText[citation.genre.toLowerCase()]}</p>}
      {citation.item_title && <h2 className="title">{citation.item_title}</h2>}
      <p>
        {citation.author}
        {citation.author && citation.date && (<span>. </span>)}
        {citation.date}
      </p>
      {renderCitation(citation)}
      <div className="citation-info">
        {citation.issn && (
          <p>
            <span className="citation-info-label">ISSN:</span> {citation.issn}
          </p>
        )}
        {citation.isbn && (
          <p>
            <span className="citation-info-label">ISBN:</span> {citation.isbn}
          </p>
        )}
        {citation.doi && (
          <p>
            <span className="citation-info-label">DOI:</span> {citation.doi}
          </p>
        )}
        {citation.pub && (
          <p>
            <span className="citation-info-label">Publisher:</span> {citation.pub}
          </p>
        )}
      </div>
    </div>
  );
};

export default Citation;
