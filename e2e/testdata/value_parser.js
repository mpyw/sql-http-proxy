// Parse values with custom logic
if (value === '' || value === 'null') return null;
if (value === 'true') return true;
if (value === 'false') return false;
if (/^\d+$/.test(value)) return parseInt(value, 10);
return value.toUpperCase();  // Custom: uppercase strings
