using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using SqlCommenterExample.Models;

namespace SqlCommenterExample.Controllers
{
    [ApiController]
    [Route("contacts")]
    public class ContactsController
    {
        private readonly SampleDbContext _db;

        public ContactsController(SampleDbContext db)
        {
            _db = db;
        }

        [HttpGet]
        public Task<List<Contact>> List()
        {
            return _db.Contacts.ToListAsync();
        }
    }
}
