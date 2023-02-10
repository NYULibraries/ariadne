//import PropTypes from 'prop-types';

const Citation = () => {

  const params = (new URL(document.location)).searchParams;

  const citation = {
    article_title: params.get("rft.atitle") || params.get("atitle"),
    journal_title: params.get("rft.jtitle") || params.get("jtitle"),
    volume: params.get("rft.volume") || params.get("volume"),
    issue: params.get("rft.issue") || params.get("issue"),
    start_page: params.get("rft.spage") || params.get("spage"),
    end_page: params.get("rft.epage") || params.get("epage"),
    genre: params.get("rft.genre") || params.get("genre"),
    issn: params.get("rft.issn") || params.get("issn"),
    date: params.get("rft.date") || params.get("date"),
    author: (params.get("rft.aulast") || params.get("aulast")) + ", " + (params.get("rft.aufirst") || params.get("aufirst"))
  };
  
  //const citation = metadataPlaceholders;

  const renderCitation = (citation) => {
    if (citation.journal_title || citation.volume || citation.issue || citation.start_page || citation.end_page) {
      return (
        <p style={{ margin: '0 0 10px' }}>
          <span style={{ boxSizing: 'border-box' }}>{citation.journal_title && 'Published in Journal'}</span>
          <span style={{ fontStyle: 'italic' }}>{citation.journal_title && citation.journal_title + '.'}</span>
          {citation.volume && 'Volume ' + citation.volume + '.'}
          {citation.issue && 'Issue ' + citation.issue + '.'}
          {citation.start_page && 'Page ' + citation.start_page}
          {citation.end_page && '-' + citation.end_page + '.'}
        </p>
      );
    }
    return null;
  };

  return (
    <div>
      {citation.genre && <p className="resource-type">{citation.genre}</p>}
      {citation.article_title && <h2 className="title">{citation.article_title}</h2>}
      <p>
        {citation.author}
        {citation.author && citation.date && ( <span>•</span>)}
        {citation.date}
      </p>
      {renderCitation(citation)}
      {citation.issn && (
        <dl className="citation-info">
          <dt>ISSN:</dt>
          <dd>{citation.issn}</dd>
        </dl>
      )}
    </div>
  );
};

export default Citation;
