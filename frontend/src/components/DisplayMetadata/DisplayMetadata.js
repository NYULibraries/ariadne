import PropTypes from 'prop-types';

const DisplayMetadata = ({ metadataPlaceholders }) => {
  return (
    <div>
      {metadataPlaceholders.genre && <p className="resource-type">{metadataPlaceholders.genre}</p>}
      {metadataPlaceholders.article_title && <h2 className="title">{metadataPlaceholders.article_title}</h2>}
      {metadataPlaceholders.author && metadataPlaceholders.date && (
        <p>
          {metadataPlaceholders.author} <span>•</span> <span aria-label="Publication date">{metadataPlaceholders.date}</span>
        </p>
      )}
      {(metadataPlaceholders.journal_title ||
        metadataPlaceholders.volume ||
        metadataPlaceholders.issue ||
        metadataPlaceholders.start_page ||
        metadataPlaceholders.end_page) && (
          <p style={{ margin: '0 0 10px' }}>
            <span style={{ boxSizing: 'border-box' }}>
              {metadataPlaceholders.journal_title && 'Published in Journal'}
            </span>
            <span style={{ fontStyle: 'italic' }}>
              {metadataPlaceholders.journal_title && metadataPlaceholders.journal_title + '.'}
            </span>
            {metadataPlaceholders.volume && 'Volume ' + metadataPlaceholders.volume + '.'}
            {metadataPlaceholders.issue && 'Issue ' + metadataPlaceholders.issue + '.'}
            {metadataPlaceholders.start_page && 'Page ' + metadataPlaceholders.start_page}
            {metadataPlaceholders.end_page && '-' + metadataPlaceholders.end_page + '.'}
          </p>
        )}
      {metadataPlaceholders.issn && (
        <dl className="citation-info">
          <dt>ISSN:</dt>
          <dd title="International Standard Serial Number">{metadataPlaceholders.issn}</dd>
        </dl>
      )}
    </div>
  );
};

DisplayMetadata.propTypes = {
  metadataPlaceholders: PropTypes.object.isRequired,
};

export default DisplayMetadata;
