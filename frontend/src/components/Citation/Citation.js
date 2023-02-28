import { useEffect } from 'react';

const Citation = () => {

  const params = (new URL(document.location)).searchParams;

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
  };
  
  // author either specified as "au", a series of separate params ("aufirst", "aulast", "auinit", "auinit1", auinitm"), 
  // or "aucorp"
  // thanks to umlaut: https://github.com/NYULibraries/umlaut/blob/master/app/models/referent.rb#L320-L340
  const getAuthorDisplayText = () => {
    var author;
    if ((author = getOpenUrlParam("au"))) {
      return author;
    } else if ((author = getOpenUrlParam("aulast"))) {
      var aufirst;
      if ((aufirst = getOpenUrlParam("aufirst"))) {
        return author + ", " + aufirst;
      } else {
        var auinit;
        if ((auinit = getOpenUrlParam("auinit"))) {
          return author + ", " + auinit;
        } else {
          if ((auinit = getOpenUrlParam("auinit1")))
            author += ", " + auinit;
          if ((auinit = getOpenUrlParam("auinitm")))
            author += auinit;
          return author;
        }
      }
    } else if ((author = getOpenUrlParam("aucorp"))) {
      return author;
    }
  }

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
    document.title = 'GetIt | ' + citation.item_title;
  }, [citation.item_title]);
  
  const renderCitation = (citation) => {
    if (citation.container_title || citation.volume || citation.issue || citation.start_page || citation.end_page) {
      return (
        <p style={{ margin: '0 0 10px' }}>
          <span style={{ boxSizing: 'border-box' }}>{citation.container_title && 'Published in '}</span>
          <span style={{ fontStyle: 'italic' }}>{citation.container_title && citation.container_title + '. '}</span>
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
    <div>
      {citation.genre && <p className="resource-type">{genresDisplayText[citation.genre.toLowerCase()]}</p>}
      {citation.item_title && <h2 className="title">{citation.item_title}</h2>}
      <p>
        {citation.author}
        {citation.author && citation.date && ( <span>•</span>)}
        {citation.date}
      </p>
      {renderCitation(citation)}
        <dl className="citation-info">
          {citation.issn && (<dt>ISSN:</dt>)}
          {citation.issn && <dd>{citation.issn}</dd>}
          {citation.isbn && (<dt>ISBN:</dt>)}
          {citation.isbn && <dd>{citation.isbn}</dd>}
          {citation.pub && (<dt>Publisher:</dt>)}
          {citation.pub && <dd>{citation.pub}</dd>}
        </dl>
    </div>
  );
};

export default Citation;
