using System.Linq;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Moq;
using TwoMQTT.Utils;
using Unifi.Liasons;
using Unifi.Models.Options;
using Unifi.Models.Shared;

namespace UnifiTest.Liasons
{
    [TestClass]
    public class MQTTLiasonTest
    {
        [TestMethod]
        public void MapDataTest()
        {
            var tests = new[] {
                new {
                    Q = new SlugMapping { MACAddress = BasicMACAddress, Slug = BasicSlug },
                    Resource = new Resource { Mac = BasicMACAddress, State = BasicState },
                    Expected = new { MACAddress = BasicMACAddress, State = BasicStateString, Slug = BasicSlug, Found = true }
                },
                new {
                    Q = new SlugMapping { MACAddress = BasicMACAddress, Slug = BasicSlug },
                    Resource = new Resource { Mac = $"{BasicMACAddress}-fake" , State = BasicState },
                    Expected = new { MACAddress = string.Empty, State = string.Empty, Slug = string.Empty, Found = false }
                },
            };

            foreach (var test in tests)
            {
                var logger = new Mock<ILogger<MQTTLiason>>();
                var generator = new Mock<IMQTTGenerator>();
                var sharedOpts = Options.Create(new SharedOpts
                {
                    Resources = new[] { test.Q }.ToList(),
                });

                generator.Setup(x => x.BuildDiscovery(It.IsAny<string>(), It.IsAny<string>(), It.IsAny<System.Reflection.AssemblyName>(), false))
                    .Returns(new TwoMQTT.Models.MQTTDiscovery());
                generator.Setup(x => x.StateTopic(test.Q.Slug, It.IsAny<string>()))
                    .Returns($"totes/{test.Q.Slug}/topic/{nameof(Resource.State)}");
                generator.Setup(x => x.BooleanOnOff(BasicState)).Returns(BasicStateString);

                var mqttLiason = new MQTTLiason(logger.Object, generator.Object, sharedOpts);
                var results = mqttLiason.MapData(test.Resource);
                var actual = results.FirstOrDefault();

                Assert.AreEqual(test.Expected.Found, results.Any(), "The mapping should exist if found.");
                if (test.Expected.Found)
                {
                    Assert.IsTrue(actual.topic.Contains(test.Expected.Slug), "The topic should contain the expected MACAddress.");
                    Assert.AreEqual(test.Expected.State, actual.payload, "The payload be the expected State.");
                }
            }
        }

        [TestMethod]
        public void DiscoveriesTest()
        {
            var tests = new[] {
                new {
                    Q = new SlugMapping { MACAddress = BasicMACAddress, Slug = BasicSlug },
                    Resource = new Resource { Mac = BasicMACAddress, State = BasicState },
                    Expected = new { MACAddress = BasicMACAddress, State = BasicState, Slug = BasicSlug }
                },
            };

            foreach (var test in tests)
            {
                var logger = new Mock<ILogger<MQTTLiason>>();
                var generator = new Mock<IMQTTGenerator>();
                var sharedOpts = Options.Create(new SharedOpts
                {
                    Resources = new[] { test.Q }.ToList(),
                });

                generator.Setup(x => x.BuildDiscovery(test.Q.Slug, It.IsAny<string>(), It.IsAny<System.Reflection.AssemblyName>(), false))
                    .Returns(new TwoMQTT.Models.MQTTDiscovery());

                var mqttLiason = new MQTTLiason(logger.Object, generator.Object, sharedOpts);
                var results = mqttLiason.Discoveries();
                var result = results.FirstOrDefault();

                Assert.IsNotNull(result, "A discovery should exist.");
            }
        }

        private static string BasicSlug = "totallyaslug";
        private static bool BasicState = true;
        private static string BasicStateString = "ON";
        private static string BasicMACAddress = "AA:BB:CC:DD:EE:FF";
    }
}
