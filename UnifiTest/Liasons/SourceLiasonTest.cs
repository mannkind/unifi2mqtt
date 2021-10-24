using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using Unifi.DataAccess;
using Unifi.Liasons;
using Unifi.Models.Options;
using Unifi.Models.Shared;

namespace UnifiTest.Liasons;

[TestClass]
public class SourceLiasonTest
{
    [TestMethod]
    public async Task FetchAllAsyncTest()
    {
        var tests = new[] {
                new {
                    Q = new SlugMapping { MACAddress = BasicMACAddress, Slug = BasicSlug },
                    Expected = new { MACAddress = BasicMACAddress, State = BasicState }
                },
            };

        foreach (var test in tests)
        {
            var logger = new Mock<ILogger<SourceLiason>>();
            var sourceDAO = new Mock<ISourceDAO>();
            var opts = Options.Create(new SourceOpts());
            var sharedOpts = Options.Create(new SharedOpts
            {
                Resources = new[] { test.Q }.ToList(),
            });

            sourceDAO.Setup(x => x.FetchOneAsync(test.Q, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(new Unifi.Models.Source.Response
                 {
                     MACAddress = test.Expected.MACAddress,
                     State = test.Expected.State,
                 });

            var sourceLiason = new SourceLiason(logger.Object, sourceDAO.Object, opts, sharedOpts);
            await foreach (var result in sourceLiason.ReceiveDataAsync())
            {
                Assert.AreEqual(test.Expected.MACAddress, result.Mac);
                Assert.AreEqual(test.Expected.State, result.State);
            }
        }
    }

    private static string BasicSlug = "totallyaslug";
    private static bool BasicState = true;
    private static string BasicMACAddress = "AA:BB:CC:DD:EE:FF";
}
